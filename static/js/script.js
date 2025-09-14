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
    loadDashboard();
    loadMahasiswa();
    loadAlumni();
    loadPekerjaan();
});

// Navigation Functions
function showSection(sectionName) {
    // Hide all sections
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });
    
    // Show selected section
    document.getElementById(sectionName).classList.add('active');
    
    // Update navbar
    document.querySelectorAll('.nav-link').forEach(link => {
        link.classList.remove('active');
    });
    
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
            
            if (data && data.length > 0) {
                data.forEach(item => {
                    tbody.innerHTML += `
                        <tr>
                            <td>${item.id}</td>
                            <td>${item.nim}</td>
                            <td>${item.nama}</td>
                            <td>${item.jurusan}</td>
                            <td>${item.angkatan}</td>
                            <td>${item.email}</td>
                            <td>
                                <button class="btn btn-sm btn-warning" onclick="editMahasiswa(${item.id})">
                                    <i class="fas fa-edit"></i>
                                </button>
                                <button class="btn btn-sm btn-danger" onclick="deleteMahasiswa(${item.id})">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbody.innerHTML = '<tr><td colspan="7" class="text-center">Tidak ada data</td></tr>';
            }
        })
        .catch(error => {
            console.error('Error:', error);
            showAlert('Error loading mahasiswa data', 'danger');
        });
}

function showMahasiswaForm(id = null) {
    currentEditId = id;
    currentEditType = 'mahasiswa';
    
    if (id) {
        // Edit mode
        fetch(`/mahasiswa/${id}`)
            .then(response => response.json())
            .then(data => {
                document.getElementById('mahasiswaId').value = data.id;
                document.getElementById('mahasiswaNim').value = data.nim;
                document.getElementById('mahasiswaNama').value = data.nama;
                document.getElementById('mahasiswaJurusan').value = data.jurusan;
                document.getElementById('mahasiswaAngkatan').value = data.angkatan;
                document.getElementById('mahasiswaEmail').value = data.email;
            });
    } else {
        // Add mode
        document.getElementById('mahasiswaForm').reset();
        document.getElementById('mahasiswaId').value = '';
    }
    
    new bootstrap.Modal(document.getElementById('mahasiswaModal')).show();
}

function saveMahasiswa() {
    const id = document.getElementById('mahasiswaId').value;
    const data = {
        nim: document.getElementById('mahasiswaNim').value,
        nama: document.getElementById('mahasiswaNama').value,
        jurusan: document.getElementById('mahasiswaJurusan').value,
        angkatan: parseInt(document.getElementById('mahasiswaAngkatan').value),
        email: document.getElementById('mahasiswaEmail').value
    };
    
    const url = id ? `/mahasiswa/${id}` : '/mahasiswa';
    const method = id ? 'PUT' : 'POST';
    
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(result => {
        showAlert('Data mahasiswa berhasil disimpan', 'success');
        bootstrap.Modal.getInstance(document.getElementById('mahasiswaModal')).hide();
        loadMahasiswa();
        loadDashboard();
    })
    .catch(error => {
        console.error('Error:', error);
        showAlert('Error saving mahasiswa data', 'danger');
    });
}

function editMahasiswa(id) {
    showMahasiswaForm(id);
}

function deleteMahasiswa(id) {
    if (confirm('Yakin ingin menghapus data mahasiswa ini?')) {
        fetch(`/mahasiswa/${id}`, {
            method: 'DELETE'
        })
        .then(() => {
            showAlert('Data mahasiswa berhasil dihapus', 'success');
            loadMahasiswa();
            loadDashboard();
        })
        .catch(error => {
            console.error('Error:', error);
            showAlert('Error deleting mahasiswa data', 'danger');
        });
    }
}

// Alumni Functions
function loadAlumni() {
    fetch('/alumni')
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById('alumniTableBody');
            tbody.innerHTML = '';
            
            if (data && data.length > 0) {
                data.forEach(item => {
                    tbody.innerHTML += `
                        <tr>
                            <td>${item.id}</td>
                            <td>${item.nim}</td>
                            <td>${item.nama}</td>
                            <td>${item.jurusan}</td>
                            <td>${item.angkatan}</td>
                            <td>${item.tahun_lulus}</td>
                            <td>${item.email}</td>
                            <td>
                                <button class="btn btn-sm btn-warning" onclick="editAlumni(${item.id})">
                                    <i class="fas fa-edit"></i>
                                </button>
                                <button class="btn btn-sm btn-danger" onclick="deleteAlumni(${item.id})">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbody.innerHTML = '<tr><td colspan="8" class="text-center">Tidak ada data</td></tr>';
            }
        })
        .catch(error => {
            console.error('Error:', error);
            showAlert('Error loading alumni data', 'danger');
        });
}

function showAlumniForm(id = null) {
    currentEditId = id;
    currentEditType = 'alumni';
    
    if (id) {
        // Edit mode
        fetch(`/alumni/${id}`)
            .then(response => response.json())
            .then(data => {
                document.getElementById('alumniId').value = data.id;
                document.getElementById('alumniNim').value = data.nim;
                document.getElementById('alumniNama').value = data.nama;
                document.getElementById('alumniJurusan').value = data.jurusan;
                document.getElementById('alumniAngkatan').value = data.angkatan;
                document.getElementById('alumniTahunLulus').value = data.tahun_lulus;
                document.getElementById('alumniEmail').value = data.email;
                document.getElementById('alumniNoTelepon').value = data.no_telepon || '';
                document.getElementById('alumniAlamat').value = data.alamat || '';
            });
    } else {
        // Add mode
        document.getElementById('alumniForm').reset();
        document.getElementById('alumniId').value = '';
    }
    
    new bootstrap.Modal(document.getElementById('alumniModal')).show();
}

function saveAlumni() {
    const id = document.getElementById('alumniId').value;
    const data = {
        nim: document.getElementById('alumniNim').value,
        nama: document.getElementById('alumniNama').value,
        jurusan: document.getElementById('alumniJurusan').value,
        angkatan: parseInt(document.getElementById('alumniAngkatan').value),
        tahun_lulus: parseInt(document.getElementById('alumniTahunLulus').value),
        email: document.getElementById('alumniEmail').value,
        no_telepon: document.getElementById('alumniNoTelepon').value,
        alamat: document.getElementById('alumniAlamat').value
    };
    
    const url = id ? `/alumni/${id}` : '/alumni';
    const method = id ? 'PUT' : 'POST';
    
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(result => {
        showAlert('Data alumni berhasil disimpan', 'success');
        bootstrap.Modal.getInstance(document.getElementById('alumniModal')).hide();
        loadAlumni();
        loadDashboard();
    })
    .catch(error => {
        console.error('Error:', error);
        showAlert('Error saving alumni data', 'danger');
    });
}

function editAlumni(id) {
    showAlumniForm(id);
}

function deleteAlumni(id) {
    if (confirm('Yakin ingin menghapus data alumni ini?')) {
        fetch(`/alumni/${id}`, {
            method: 'DELETE'
        })
        .then(() => {
            showAlert('Data alumni berhasil dihapus', 'success');
            loadAlumni();
            loadDashboard();
        })
        .catch(error => {
            console.error('Error:', error);
            showAlert('Error deleting alumni data', 'danger');
        });
    }
}

// Pekerjaan Functions
function loadPekerjaan() {
    fetch('/pekerjaan')
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById('pekerjaanTableBody');
            tbody.innerHTML = '';
            
            if (data && data.length > 0) {
                data.forEach(item => {
                    tbody.innerHTML += `
                        <tr>
                            <td>${item.id}</td>
                            <td>Alumni ID: ${item.alumni_id}</td>
                            <td>${item.nama_perusahaan}</td>
                            <td>${item.posisi_jabatan}</td>
                            <td>${item.bidang_industri}</td>
                            <td>${item.lokasi_kerja}</td>
                            <td><span class="badge bg-${getStatusColor(item.status_pekerjaan)}">${item.status_pekerjaan}</span></td>
                            <td>
                                <button class="btn btn-sm btn-warning" onclick="editPekerjaan(${item.id})">
                                    <i class="fas fa-edit"></i>
                                </button>
                                <button class="btn btn-sm btn-danger" onclick="deletePekerjaan(${item.id})">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbody.innerHTML = '<tr><td colspan="8" class="text-center">Tidak ada data</td></tr>';
            }
        })
        .catch(error => {
            console.error('Error:', error);
            showAlert('Error loading pekerjaan data', 'danger');
        });
}

function loadAlumniOptions() {
    fetch('/alumni')
        .then(response => response.json())
        .then(data => {
            const select = document.getElementById('pekerjaanAlumniId');
            select.innerHTML = '<option value="">Pilih Alumni</option>';
            
            if (data && data.length > 0) {
                data.forEach(alumni => {
                    select.innerHTML += `<option value="${alumni.id}">${alumni.nama} (${alumni.nim})</option>`;
                });
            }
        })
        .catch(error => {
            console.error('Error loading alumni options:', error);
        });
}

function showPekerjaanForm(id = null) {
    currentEditId = id;
    currentEditType = 'pekerjaan';
    loadAlumniOptions();
    
    if (id) {
        // Edit mode
        fetch(`/pekerjaan/${id}`)
            .then(response => response.json())
            .then(data => {
                document.getElementById('pekerjaanId').value = data.id;
                document.getElementById('pekerjaanAlumniId').value = data.alumni_id;
                document.getElementById('pekerjaanNamaPerusahaan').value = data.nama_perusahaan;
                document.getElementById('pekerjaanPosisiJabatan').value = data.posisi_jabatan;
                document.getElementById('pekerjaanBidangIndustri').value = data.bidang_industri;
                document.getElementById('pekerjaanLokasiKerja').value = data.lokasi_kerja;
                document.getElementById('pekerjaanGajiRange').value = data.gaji_range || '';
                document.getElementById('pekerjaanTanggalMulai').value = data.tanggal_mulai_kerja ? data.tanggal_mulai_kerja.split('T')[0] : '';
                document.getElementById('pekerjaanTanggalSelesai').value = data.tanggal_selesai_kerja ? data.tanggal_selesai_kerja.split('T')[0] : '';
                document.getElementById('pekerjaanStatus').value = data.status_pekerjaan;
                document.getElementById('pekerjaanDeskripsi').value = data.deskripsi_pekerjaan || '';
            });
    } else {
        // Add mode
        document.getElementById('pekerjaanForm').reset();
        document.getElementById('pekerjaanId').value = '';
    }
    
    new bootstrap.Modal(document.getElementById('pekerjaanModal')).show();
}

function savePekerjaan() {
    const id = document.getElementById('pekerjaanId').value;
    const data = {
        alumni_id: parseInt(document.getElementById('pekerjaanAlumniId').value),
        nama_perusahaan: document.getElementById('pekerjaanNamaPerusahaan').value,
        posisi_jabatan: document.getElementById('pekerjaanPosisiJabatan').value,
        bidang_industri: document.getElementById('pekerjaanBidangIndustri').value,
        lokasi_kerja: document.getElementById('pekerjaanLokasiKerja').value,
        gaji_range: document.getElementById('pekerjaanGajiRange').value,
        tanggal_mulai_kerja: document.getElementById('pekerjaanTanggalMulai').value,
        tanggal_selesai_kerja: document.getElementById('pekerjaanTanggalSelesai').value || null,
        status_pekerjaan: document.getElementById('pekerjaanStatus').value,
        deskripsi_pekerjaan: document.getElementById('pekerjaanDeskripsi').value
    };
    
    const url = id ? `/pekerjaan/${id}` : '/pekerjaan';
    const method = id ? 'PUT' : 'POST';
    
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(result => {
        showAlert('Data pekerjaan berhasil disimpan', 'success');
        bootstrap.Modal.getInstance(document.getElementById('pekerjaanModal')).hide();
        loadPekerjaan();
        loadDashboard();
    })
    .catch(error => {
        console.error('Error:', error);
        showAlert('Error saving pekerjaan data', 'danger');
    });
}

function editPekerjaan(id) {
    showPekerjaanForm(id);
}

function deletePekerjaan(id) {
    if (confirm('Yakin ingin menghapus data pekerjaan ini?')) {
        fetch(`/pekerjaan/${id}`, {
            method: 'DELETE'
        })
        .then(() => {
            showAlert('Data pekerjaan berhasil dihapus', 'success');
            loadPekerjaan();
            loadDashboard();
        })
        .catch(error => {
            console.error('Error:', error);
            showAlert('Error deleting pekerjaan data', 'danger');
        });
    }
}

// Utility Functions
function getStatusColor(status) {
    switch(status) {
        case 'aktif': return 'success';
        case 'selesai': return 'primary';
        case 'resigned': return 'danger';
        default: return 'secondary';
    }
}

function showAlert(message, type) {
    // Create alert element
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    // Insert at top of container
    const container = document.querySelector('.container');
    container.insertBefore(alertDiv, container.firstChild);
    
    // Auto-dismiss after 3 seconds
    setTimeout(() => {
        if (alertDiv) {
            alertDiv.remove();
        }
    }, 3000);
}
