// Global Variables
let currentEditId = null;
let currentEditType = null;

// Helper function untuk authorized requests
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login';
        return null;
    }
    return {
        'Authorization': 'Bearer ' + token,
        'Content-Type': 'application/json'
    };
}

function authorizedFetch(url, options = {}) {
    const headers = getAuthHeaders();
    if (!headers) return Promise.reject('No auth token');
    
    return fetch(url, {
        ...options,
        headers: {
            ...headers,
            ...(options.headers || {})
        }
    }).then(response => {
        if (response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
            throw new Error('Unauthorized');
        }
        return response;
    });
}

// Initialize App
document.addEventListener('DOMContentLoaded', function() {
    // Check authentication first
    checkAuth().then(() => {
        loadDashboard();
        loadMahasiswa();
        loadAlumni();
        loadPekerjaan();
    });
});

// Check authentication
function checkAuth() {
    return authorizedFetch('/api/profile')
        .then(response => response.json())
        .then(data => {
            console.log('User authenticated:', data.user);
            showAuthenticatedState(data.user);
            return data;
        })
        .catch(error => {
            console.error('Auth check failed:', error);
            window.location.href = '/login';
            throw error;
        });
}

// Show authenticated state in UI
function showAuthenticatedState(user) {
    // Hide auth buttons, show user info and logout
    const authButtons = document.getElementById('authButtons');
    const userInfo = document.getElementById('userInfo');
    const logoutBtn = document.getElementById('logoutBtn');
    const userName = document.getElementById('userName');
    const userRole = document.getElementById('userRole');
    
    if (authButtons) authButtons.style.display = 'none';
    if (userInfo) userInfo.style.display = 'inline';
    if (logoutBtn) logoutBtn.style.display = 'inline-block';
    if (userName) userName.textContent = user.username;
    if (userRole) userRole.textContent = user.role;
    
    // Show protected content
    document.querySelectorAll('.protected-content').forEach(element => {
        element.style.display = 'block';
    });
    
    // Show admin-only content if user is admin
    if (user.role === 'admin') {
        document.querySelectorAll('.admin-only').forEach(element => {
            element.style.display = 'block';
        });
    }
}

// Check if current user is admin
function isAdmin() {
    const userRoleElement = document.getElementById('userRole');
    return userRoleElement && userRoleElement.textContent === 'admin';
}

// Logout function
function logout() {
    const token = localStorage.getItem('token');
    
    // Call logout endpoint
    if (token) {
        fetch('/api/logout', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        }).catch(error => {
            console.log('Logout API call failed:', error);
        });
    }
    
    // Remove token
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    // Redirect to login
    window.location.href = '/login';
}

// Navigation Functions
function showSection(sectionName) {
    // Hide all sections
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });
    
    // Show selected section
    document.getElementById(sectionName).classList.add('active');
    
    // Load data for the section
    switch(sectionName) {
        case 'mahasiswa':
            loadMahasiswa();
            break;
        case 'alumni':
            loadAlumni();
            break;
        case 'pekerjaan':
            loadPekerjaan();
            loadAlumniOptions();
            break;
        default:
            loadDashboard();
    }
}

// Dashboard Functions
function loadDashboard() {
    // Load counts for dashboard using authorized fetch
    authorizedFetch('/api/mahasiswa/count')
        .then(response => response.json())
        .then(data => {
            document.getElementById('totalMahasiswa').textContent = data.count || 0;
        })
        .catch(error => {
            console.error('Error loading mahasiswa count:', error);
            document.getElementById('totalMahasiswa').textContent = 'Error';
        });

    authorizedFetch('/api/alumni/count')
        .then(response => response.json())
        .then(data => {
            document.getElementById('totalAlumni').textContent = data.count || 0;
        })
        .catch(error => {
            console.error('Error loading alumni count:', error);
            document.getElementById('totalAlumni').textContent = 'Error';
        });

    authorizedFetch('/api/pekerjaan/count')
        .then(response => response.json())
        .then(data => {
            document.getElementById('totalPekerjaan').textContent = data.count || 0;
        })
        .catch(error => {
            console.error('Error loading pekerjaan count:', error);
            document.getElementById('totalPekerjaan').textContent = 'Error';
        });
}

// Mahasiswa Functions
function loadMahasiswa() {
    authorizedFetch('/api/mahasiswa')
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById('mahasiswaTableBody');
            tbody.innerHTML = '';
            
            data.forEach(mahasiswa => {
                const row = tbody.insertRow();
                if (isAdmin()) {
                    row.innerHTML = `
                        <td>${mahasiswa.id}</td>
                        <td>${mahasiswa.nim}</td>
                        <td>${mahasiswa.nama}</td>
                        <td>${mahasiswa.email}</td>
                        <td>${mahasiswa.jurusan}</td>
                        <td>${mahasiswa.angkatan}</td>
                        <td>
                            <button class="btn btn-sm btn-warning" onclick="editMahasiswa(${mahasiswa.id})">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deleteMahasiswa(${mahasiswa.id})">
                                <i class="fas fa-trash"></i>
                            </button>
                        </td>
                    `;
                } else {
                    row.innerHTML = `
                        <td>${mahasiswa.id}</td>
                        <td>${mahasiswa.nim}</td>
                        <td>${mahasiswa.nama}</td>
                        <td>${mahasiswa.email}</td>
                        <td>${mahasiswa.jurusan}</td>
                        <td>${mahasiswa.angkatan}</td>
                    `;
                }
            });
        })
        .catch(error => {
            console.error('Error loading mahasiswa:', error);
            showAlert('danger', 'Gagal memuat data mahasiswa');
        });
}

function showMahasiswaForm() {
    clearMahasiswaForm();
    currentEditId = null;
    currentEditType = 'mahasiswa';
    const modal = new bootstrap.Modal(document.getElementById('mahasiswaModal'));
    modal.show();
}

function editMahasiswa(id) {
    authorizedFetch(`/api/mahasiswa/${id}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById('mahasiswaId').value = data.id;
            document.getElementById('mahasiswaNim').value = data.nim;
            document.getElementById('mahasiswaNama').value = data.nama;
            document.getElementById('mahasiswaEmail').value = data.email;
            document.getElementById('mahasiswaJurusan').value = data.jurusan;
            document.getElementById('mahasiswaAngkatan').value = data.angkatan;
            
            currentEditId = id;
            currentEditType = 'mahasiswa';
            
            const modal = new bootstrap.Modal(document.getElementById('mahasiswaModal'));
            modal.show();
        })
        .catch(error => {
            console.error('Error loading mahasiswa:', error);
            showAlert('danger', 'Gagal memuat data mahasiswa');
        });
}

function saveMahasiswa() {
    const id = document.getElementById('mahasiswaId').value;
    const mahasiswaData = {
        nim: document.getElementById('mahasiswaNim').value,
        nama: document.getElementById('mahasiswaNama').value,
        email: document.getElementById('mahasiswaEmail').value,
        jurusan: document.getElementById('mahasiswaJurusan').value,
        angkatan: parseInt(document.getElementById('mahasiswaAngkatan').value)
    };

    const url = id ? `/api/mahasiswa/${id}` : '/api/mahasiswa';
    const method = id ? 'PUT' : 'POST';

    authorizedFetch(url, {
        method: method,
        body: JSON.stringify(mahasiswaData)
    })
    .then(response => response.json())
    .then(data => {
        showAlert('success', `Mahasiswa berhasil ${id ? 'diupdate' : 'ditambahkan'}`);
        const modal = bootstrap.Modal.getInstance(document.getElementById('mahasiswaModal'));
        modal.hide();
        loadMahasiswa();
        clearMahasiswaForm();
    })
    .catch(error => {
        console.error('Error saving mahasiswa:', error);
        showAlert('danger', `Gagal ${id ? 'mengupdate' : 'menambahkan'} mahasiswa`);
    });
}

function deleteMahasiswa(id) {
    if (confirm('Apakah Anda yakin ingin menghapus mahasiswa ini?')) {
        authorizedFetch(`/api/mahasiswa/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                showAlert('success', 'Mahasiswa berhasil dihapus');
                loadMahasiswa();
            } else {
                throw new Error('Delete failed');
            }
        })
        .catch(error => {
            console.error('Error deleting mahasiswa:', error);
            showAlert('danger', 'Gagal menghapus mahasiswa');
        });
    }
}

// Alumni Functions
function loadAlumni() {
    authorizedFetch('/api/alumni')
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById('alumniTableBody');
            tbody.innerHTML = '';
            
            data.forEach(alumni => {
                const row = tbody.insertRow();
                if (isAdmin()) {
                    row.innerHTML = `
                        <td>${alumni.id}</td>
                        <td>${alumni.nim}</td>
                        <td>${alumni.nama}</td>
                        <td>${alumni.jurusan}</td>
                        <td>${alumni.tahun_lulus}</td>
                        <td>${alumni.email}</td>
                        <td>
                            <button class="btn btn-sm btn-warning" onclick="editAlumni(${alumni.id})">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deleteAlumni(${alumni.id})">
                                <i class="fas fa-trash"></i>
                            </button>
                        </td>
                    `;
                } else {
                    row.innerHTML = `
                        <td>${alumni.id}</td>
                        <td>${alumni.nim}</td>
                        <td>${alumni.nama}</td>
                        <td>${alumni.jurusan}</td>
                        <td>${alumni.tahun_lulus}</td>
                        <td>${alumni.email}</td>
                    `;
                }
            });
        })
        .catch(error => {
            console.error('Error loading alumni:', error);
            showAlert('danger', 'Gagal memuat data alumni');
        });
}

function showAlumniForm() {
    clearAlumniForm();
    currentEditId = null;
    currentEditType = 'alumni';
    const modal = new bootstrap.Modal(document.getElementById('alumniModal'));
    modal.show();
}

function editAlumni(id) {
    authorizedFetch(`/api/alumni/${id}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById('alumniId').value = data.id;
            document.getElementById('alumniNim').value = data.nim;
            document.getElementById('alumniNama').value = data.nama;
            document.getElementById('alumniEmail').value = data.email;
            document.getElementById('alumniJurusan').value = data.jurusan;
            document.getElementById('alumniTahunLulus').value = data.tahun_lulus;
            document.getElementById('alumniAlamat').value = data.alamat;
            
            currentEditId = id;
            currentEditType = 'alumni';
            
            const modal = new bootstrap.Modal(document.getElementById('alumniModal'));
            modal.show();
        })
        .catch(error => {
            console.error('Error loading alumni:', error);
            showAlert('danger', 'Gagal memuat data alumni');
        });
}

function saveAlumni() {
    const id = document.getElementById('alumniId').value;
    const alumniData = {
        nim: document.getElementById('alumniNim').value,
        nama: document.getElementById('alumniNama').value,
        email: document.getElementById('alumniEmail').value,
        jurusan: document.getElementById('alumniJurusan').value,
        tahun_lulus: parseInt(document.getElementById('alumniTahunLulus').value),
        alamat: document.getElementById('alumniAlamat').value
    };

    const url = id ? `/api/alumni/${id}` : '/api/alumni';
    const method = id ? 'PUT' : 'POST';

    authorizedFetch(url, {
        method: method,
        body: JSON.stringify(alumniData)
    })
    .then(response => response.json())
    .then(data => {
        showAlert('success', `Alumni berhasil ${id ? 'diupdate' : 'ditambahkan'}`);
        const modal = bootstrap.Modal.getInstance(document.getElementById('alumniModal'));
        modal.hide();
        loadAlumni();
        clearAlumniForm();
    })
    .catch(error => {
        console.error('Error saving alumni:', error);
        showAlert('danger', `Gagal ${id ? 'mengupdate' : 'menambahkan'} alumni`);
    });
}

function deleteAlumni(id) {
    if (confirm('Apakah Anda yakin ingin menghapus alumni ini?')) {
        authorizedFetch(`/api/alumni/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                showAlert('success', 'Alumni berhasil dihapus');
                loadAlumni();
            } else {
                throw new Error('Delete failed');
            }
        })
        .catch(error => {
            console.error('Error deleting alumni:', error);
            showAlert('danger', 'Gagal menghapus alumni');
        });
    }
}

// Pekerjaan Functions
function loadPekerjaan() {
    authorizedFetch('/api/pekerjaan')
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById('pekerjaanTableBody');
            tbody.innerHTML = '';
            
            data.forEach(pekerjaan => {
                const row = tbody.insertRow();
                if (isAdmin()) {
                    row.innerHTML = `
                        <td>${pekerjaan.id}</td>
                        <td>${pekerjaan.alumni.nama || 'N/A'}</td>
                        <td>${pekerjaan.nama_perusahaan}</td>
                        <td>${pekerjaan.posisi_jabatan}</td>
                        <td>${pekerjaan.bidang_industri}</td>
                        <td>${pekerjaan.lokasi_kerja}</td>
                        <td>${pekerjaan.status}</td>
                        <td>
                            <button class="btn btn-sm btn-warning" onclick="editPekerjaan(${pekerjaan.id})">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deletePekerjaan(${pekerjaan.id})">
                                <i class="fas fa-trash"></i>
                            </button>
                        </td>
                    `;
                } else {
                    row.innerHTML = `
                        <td>${pekerjaan.id}</td>
                        <td>${pekerjaan.alumni.nama || 'N/A'}</td>
                        <td>${pekerjaan.nama_perusahaan}</td>
                        <td>${pekerjaan.posisi_jabatan}</td>
                        <td>${pekerjaan.bidang_industri}</td>
                        <td>${pekerjaan.lokasi_kerja}</td>
                        <td>${pekerjaan.status}</td>
                    `;
                }
            });
        })
        .catch(error => {
            console.error('Error loading pekerjaan:', error);
            showAlert('danger', 'Gagal memuat data pekerjaan');
        });
}

function showPekerjaanForm() {
    clearPekerjaanForm();
    currentEditId = null;
    currentEditType = 'pekerjaan';
    loadAlumniOptions();
    const modal = new bootstrap.Modal(document.getElementById('pekerjaanModal'));
    modal.show();
}

function loadAlumniOptions() {
    authorizedFetch('/api/alumni')
        .then(response => response.json())
        .then(data => {
            const select = document.getElementById('pekerjaanAlumniId');
            select.innerHTML = '<option value="">Pilih Alumni</option>';
            
            data.forEach(alumni => {
                const option = document.createElement('option');
                option.value = alumni.id;
                option.textContent = `${alumni.nama} (${alumni.nim})`;
                select.appendChild(option);
            });
        })
        .catch(error => {
            console.error('Error loading alumni options:', error);
        });
}

function editPekerjaan(id) {
    authorizedFetch(`/api/pekerjaan/${id}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById('pekerjaanId').value = data.id;
            document.getElementById('pekerjaanAlumniId').value = data.alumni_id;
            document.getElementById('pekerjaanNamaPerusahaan').value = data.nama_perusahaan;
            document.getElementById('pekerjaanPosisiJabatan').value = data.posisi_jabatan;
            document.getElementById('pekerjaanBidangIndustri').value = data.bidang_industri;
            document.getElementById('pekerjaanLokasiKerja').value = data.lokasi_kerja;
            document.getElementById('pekerjaanGajiRange').value = data.gaji_range;
            document.getElementById('pekerjaanTanggalMulai').value = data.tanggal_mulai;
            document.getElementById('pekerjaanTanggalSelesai').value = data.tanggal_selesai;
            document.getElementById('pekerjaanStatus').value = data.status;
            
            currentEditId = id;
            currentEditType = 'pekerjaan';
            
            const modal = new bootstrap.Modal(document.getElementById('pekerjaanModal'));
            modal.show();
        })
        .catch(error => {
            console.error('Error loading pekerjaan:', error);
            showAlert('danger', 'Gagal memuat data pekerjaan');
        });
}

function savePekerjaan() {
    const id = document.getElementById('pekerjaanId').value;
    const pekerjaanData = {
        alumni_id: parseInt(document.getElementById('pekerjaanAlumniId').value),
        nama_perusahaan: document.getElementById('pekerjaanNamaPerusahaan').value,
        posisi_jabatan: document.getElementById('pekerjaanPosisiJabatan').value,
        bidang_industri: document.getElementById('pekerjaanBidangIndustri').value,
        lokasi_kerja: document.getElementById('pekerjaanLokasiKerja').value,
        gaji_range: document.getElementById('pekerjaanGajiRange').value,
        tanggal_mulai: document.getElementById('pekerjaanTanggalMulai').value,
        tanggal_selesai: document.getElementById('pekerjaanTanggalSelesai').value,
        status: document.getElementById('pekerjaanStatus').value
    };

    const url = id ? `/api/pekerjaan/${id}` : '/api/pekerjaan';
    const method = id ? 'PUT' : 'POST';

    authorizedFetch(url, {
        method: method,
        body: JSON.stringify(pekerjaanData)
    })
    .then(response => response.json())
    .then(data => {
        showAlert('success', `Pekerjaan berhasil ${id ? 'diupdate' : 'ditambahkan'}`);
        const modal = bootstrap.Modal.getInstance(document.getElementById('pekerjaanModal'));
        modal.hide();
        loadPekerjaan();
        clearPekerjaanForm();
    })
    .catch(error => {
        console.error('Error saving pekerjaan:', error);
        showAlert('danger', `Gagal ${id ? 'mengupdate' : 'menambahkan'} pekerjaan`);
    });
}

function deletePekerjaan(id) {
    if (confirm('Apakah Anda yakin ingin menghapus pekerjaan ini?')) {
        authorizedFetch(`/api/pekerjaan/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                showAlert('success', 'Pekerjaan berhasil dihapus');
                loadPekerjaan();
            } else {
                throw new Error('Delete failed');
            }
        })
        .catch(error => {
            console.error('Error deleting pekerjaan:', error);
            showAlert('danger', 'Gagal menghapus pekerjaan');
        });
    }
}

// Form clearing functions
function clearMahasiswaForm() {
    document.getElementById('mahasiswaForm').reset();
    document.getElementById('mahasiswaId').value = '';
    currentEditId = null;
    currentEditType = null;
}

function clearAlumniForm() {
    document.getElementById('alumniForm').reset();
    document.getElementById('alumniId').value = '';
    currentEditId = null;
    currentEditType = null;
}

function clearPekerjaanForm() {
    document.getElementById('pekerjaanForm').reset();
    document.getElementById('pekerjaanId').value = '';
    currentEditId = null;
    currentEditType = null;
}

// Utility functions
function showAlert(type, message) {
    // Create alert element
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    // Add to container
    const container = document.querySelector('.container');
    container.insertBefore(alertDiv, container.firstChild);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (alertDiv.parentNode) {
            alertDiv.remove();
        }
    }, 5000);
}

// Modal event handlers
function openMahasiswaModal() {
    clearMahasiswaForm();
    const modal = new bootstrap.Modal(document.getElementById('mahasiswaModal'));
    modal.show();
}

function openAlumniModal() {
    clearAlumniForm();
    const modal = new bootstrap.Modal(document.getElementById('alumniModal'));
    modal.show();
}

function openPekerjaanModal() {
    clearPekerjaanForm();
    loadAlumniOptions();
    const modal = new bootstrap.Modal(document.getElementById('pekerjaanModal'));
    modal.show();
}
