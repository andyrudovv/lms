import React, { useState, useEffect, createContext, useContext } from 'react';
import {
    BookOpen, Users, Calendar, BarChart3, LogOut, User, Menu, X,
    Home, GraduationCap, UserCircle, Shield, ChevronDown, Plus,
    CheckCircle, XCircle, Clock, Trash2, UserPlus, Settings
} from 'lucide-react';

const AuthContext = createContext(null);
const useAuth = () => {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error('useAuth must be used within AuthProvider');
    return ctx;
};

const API_BASE = 'http://localhost:8080/api/v1';

const api = {
    async request(endpoint, options = {}) {
        const token = localStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json',
            ...(token && { Authorization: `Bearer ${token}` }),
            ...options.headers,
        };
        const response = await fetch(`${API_BASE}${endpoint}`, { ...options, headers });
        const json = await response.json().catch(() => ({}));
        if (!response.ok) {
            const msg = json?.error?.message || json?.message || 'Request failed';
            throw new Error(msg);
        }
        return json.data !== undefined ? json.data : json;
    },
    register:             (data)           => api.request('/auth/register', { method: 'POST', body: JSON.stringify(data) }),
    login:                (data)           => api.request('/auth/login',    { method: 'POST', body: JSON.stringify(data) }),
    getProfile:           ()               => api.request('/me'),
    getCourses:           ()               => api.request('/courses').then(d => d.items ?? []),
    getMyCourses:         ()               => api.request('/my/courses').then(d => d.items ?? []),
    createCourse:         (data)           => api.request('/courses', { method: 'POST', body: JSON.stringify(data) }),
    enrollStudent:        (courseId, sid)  => api.request(`/courses/${courseId}/enroll`, { method: 'POST', body: JSON.stringify({ student_id: sid }) }),
    markAttendance:       (courseId, data) => api.request(`/courses/${courseId}/attendance`, { method: 'POST', body: JSON.stringify(data) }),
    getCourseAttendance:  (courseId)       => api.request(`/courses/${courseId}/attendance`).then(d => d.items ?? []),
    getMyAttendance:      (courseId)       => api.request(`/my/attendance?course_id=${courseId}`).then(d => d.items ?? []),
    getUsers:             ()               => api.request('/users').then(d => d.items ?? []),
    createUser:           (data)           => api.request('/users', { method: 'POST', body: JSON.stringify(data) }),
    updateUserRole:       (userId, role)   => api.request(`/users/${userId}/role`, { method: 'PATCH', body: JSON.stringify({ role }) }),
    getRoles:             ()               => api.request('/roles').then(d => d.items ?? []),
    getStudentsInCourse:   (courseId)       => api.request(`/courses/${courseId}/students`).then(d => d.items ?? []),
    getAvailableStudentsForCourse: (courseId)       => api.request(`/courses/${courseId}/available-students`).then(d => d.items ?? []),
};

// ─── Shared helpers ────────────────────────────────────────────────────────

function Modal({ title, onClose, children }) {
    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal" onClick={e => e.stopPropagation()}>
                <div className="modal-header">
                    <h2>{title}</h2>
                    <button onClick={onClose}><X size={20} /></button>
                </div>
                {children}
            </div>
        </div>
    );
}

function Select({ value, onChange, options, placeholder }) {
    return (
        <div className="select-wrapper">
            <select value={value} onChange={e => onChange(e.target.value)} className="form-select">
                {placeholder && <option value="">{placeholder}</option>}
                {options.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
            </select>
            <ChevronDown size={16} className="select-icon" />
        </div>
    );
}

const ROLE_COLORS = { admin: '#e74c3c', teacher: '#3b6ea8', student: '#5fa8a0' };
const ROLE_BG    = { admin: 'rgba(231,76,60,0.12)', teacher: 'rgba(59,110,168,0.12)', student: 'rgba(95,168,160,0.12)' };

function RolePill({ role }) {
    return (
        <span style={{
            display: 'inline-block', padding: '0.3rem 0.8rem', borderRadius: '20px',
            fontSize: '0.8rem', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em',
            background: ROLE_BG[role] || 'var(--bg-dark)',
            color: ROLE_COLORS[role] || 'var(--text-muted)',
        }}>{role}</span>
    );
}

// ─── Login ─────────────────────────────────────────────────────────────────

function LoginPage({ onLogin }) {
    const [isLogin, setIsLogin] = useState(true);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [form, setForm] = useState({ email: '', password: '', fullName: '' });

    const set = k => e => setForm(f => ({ ...f, [k]: e.target.value }));

    const handleSubmit = async e => {
        e.preventDefault();
        setLoading(true); setError('');
        try {
            if (isLogin) {
                const res = await api.login({ email: form.email, password: form.password });
                localStorage.setItem('token', res.access_token);
                onLogin();
            } else {
                await api.register({ email: form.email, password: form.password, full_name: form.fullName });
                setIsLogin(true);
                setError('Registration successful! Please login.');
            }
        } catch (err) { setError(err.message); }
        finally { setLoading(false); }
    };

    return (
        <div className="login-container">
            <div className="login-background" />
            <div className="login-card">
                <div className="login-header">
                    <GraduationCap size={48} />
                    <h1>AITU LMS</h1>
                    <p>Learning Management System</p>
                </div>
                <div className="auth-tabs">
                    <button className={isLogin ? 'active' : ''} onClick={() => setIsLogin(true)}>Login</button>
                    <button className={!isLogin ? 'active' : ''} onClick={() => setIsLogin(false)}>Register</button>
                </div>
                <form onSubmit={handleSubmit} className="login-form">
                    {!isLogin && (
                        <div className="form-group">
                            <label>Full Name</label>
                            <input type="text" value={form.fullName} onChange={set('fullName')} required placeholder="Enter your full name" />
                        </div>
                    )}
                    <div className="form-group">
                        <label>Email</label>
                        <input type="email" value={form.email} onChange={set('email')} required placeholder="your@email.com" />
                    </div>
                    <div className="form-group">
                        <label>Password</label>
                        <input type="password" value={form.password} onChange={set('password')} required placeholder="••••••••" />
                    </div>
                    {error && <div className={`alert ${error.includes('successful') ? 'success' : 'error'}`}>{error}</div>}
                    <button type="submit" className="btn-primary" disabled={loading}>
                        {loading ? 'Please wait...' : isLogin ? 'Login' : 'Register'}
                    </button>
                </form>
            </div>
        </div>
    );
}

// ─── Profile ───────────────────────────────────────────────────────────────

function ProfilePage({ user }) {
    return (
        <div className="dashboard-content">
            <h1 className="page-title">My Profile</h1>
            <div className="profile-container">
                <div className="profile-card">
                    <div className="profile-avatar"><UserCircle size={80} /></div>
                    <div className="profile-info">
                        <h2>{user.full_name}</h2>
                        <RolePill role={user.role} />
                    </div>
                </div>
                <div className="profile-details">
                    {[['Email', user.email], ['User ID', user.id], ['Role', user.role]].map(([label, val]) => (
                        <div key={label} className="detail-item">
                            <label>{label}</label>
                            <p>{val}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

// ─── STUDENT views ─────────────────────────────────────────────────────────

function StudentHome({ user }) {
    const [myCourses, setMyCourses] = useState([]);
    const [totalCourses, setTotalCourses] = useState(0);

    useEffect(() => {
        Promise.all([api.getMyCourses(), api.getCourses()]).then(([mine, all]) => {
            setMyCourses(mine); setTotalCourses(all.length);
        }).catch(console.error);
    }, []);

    return (
        <div className="dashboard-content">
            <div className="welcome-section">
                <h1>Welcome back, {user.full_name}!</h1>
                <p className="subtitle">Here's your learning overview</p>
            </div>
            <div className="stats-grid">
                <StatCard icon={<BookOpen />} value={myCourses.length} label="Enrolled Courses" />
                <StatCard icon={<Users />} value={totalCourses} label="Total Courses" />
                <StatCard icon={<Calendar />} value={new Date().toLocaleDateString('en-US', { weekday: 'long' })} label="Today" />
                <StatCard icon={<BarChart3 />} value="Active" label="Status" />
            </div>
            <div className="section">
                <h2>My Courses</h2>
                {myCourses.length === 0
                    ? <EmptyState icon={<BookOpen size={48} />} text="No courses enrolled yet" />
                    : <div className="courses-grid">{myCourses.map(c => <CourseCard key={c.id} course={c} />)}</div>}
            </div>
        </div>
    );
}

function StudentCourses() {
    const [courses, setCourses] = useState([]);
    useEffect(() => { api.getCourses().then(setCourses).catch(console.error); }, []);
    return (
        <div className="dashboard-content">
            <div className="page-header"><h1 className="page-title">All Courses</h1></div>
            {courses.length === 0
                ? <EmptyState icon={<BookOpen size={48} />} text="No courses available" />
                : <div className="courses-grid">{courses.map(c => <CourseCard key={c.id} course={c} />)}</div>}
        </div>
    );
}

function StudentAttendance({ user }) {
    const [courses, setCourses] = useState([]);
    const [selected, setSelected] = useState(null);
    const [attendance, setAttendance] = useState([]);

    useEffect(() => { api.getMyCourses().then(setCourses).catch(console.error); }, []);

    const selectCourse = async c => {
        setSelected(c);
        const data = await api.getMyAttendance(c.id).catch(() => []);
        setAttendance(data);
    };

    return (
        <div className="dashboard-content">
            <h1 className="page-title">My Attendance</h1>
            <div className="attendance-container">
                <CourseSidebar courses={courses} selected={selected} onSelect={selectCourse} />
                <div className="attendance-list">
                    {!selected
                        ? <EmptyState icon={<Calendar size={48} />} text="Select a course to view attendance" />
                        : attendance.length === 0
                            ? <EmptyState icon={<Calendar size={48} />} text="No attendance records yet" />
                            : <AttendanceTable records={attendance} showStudent={false} />}
                </div>
            </div>
        </div>
    );
}

// ─── TEACHER views ─────────────────────────────────────────────────────────

function TeacherHome({ user }) {
    const [courses, setCourses] = useState([]);
    useEffect(() => { api.getMyCourses().then(setCourses).catch(console.error); }, []);
    return (
        <div className="dashboard-content">
            <div className="welcome-section">
                <h1>Welcome back, {user.full_name}!</h1>
                <p className="subtitle">Manage your courses and track student progress</p>
            </div>
            <div className="stats-grid">
                <StatCard icon={<BookOpen />} value={courses.length} label="My Courses" />
                <StatCard icon={<Calendar />} value={new Date().toLocaleDateString('en-US', { weekday: 'long' })} label="Today" />
                <StatCard icon={<BarChart3 />} value="Active" label="Status" />
            </div>
            <div className="section">
                <h2>My Courses</h2>
                {courses.length === 0
                    ? <EmptyState icon={<BookOpen size={48} />} text="No courses yet — create one in Courses" />
                    : <div className="courses-grid">{courses.map(c => <CourseCard key={c.id} course={c} />)}</div>}
            </div>
        </div>
    );
}

function TeacherCourses({ user }) {
    const [courses, setCourses] = useState([]);
    const [usersInCourse, setUsersInCourse] = useState([]);
    const [usersForEnrollment, setUsersForEnrollment] = useState([]);
    const [showCreate, setShowCreate] = useState(false);
    const [showEnroll, setShowEnroll] = useState(null); // course object
    const [showMark, setShowMark] = useState(null);     // course object
    const [newTitle, setNewTitle] = useState('');
    const [loading, setLoading] = useState(false);

    const [currentCourseId, setCurrentCourseId] = useState(null);

    // Load teacher's courses
    const loadCourses = async () => {
        try {
            const c = await api.getMyCourses(); // should call backend endpoint for teacher's courses
            setCourses(c);
        } catch (err) {
            console.error(err);
        }
    };

    // Load students for a specific course
    const loadStudents = async (courseId) => {
        if (!courseId) return;
        try {
            const res = await api.getStudentsInCourse(courseId);
            console.log('Students in course response:', res);
            const studentsArray = Array.isArray(res) ? res : res.items ?? [];
            setUsersInCourse(studentsArray);
        } catch (err) {
            console.error(err);
        }
    };


    const loadStudentsAvailableForEnrollment = async (courseId) => {
        if (!courseId) return;
        try {            
            const res = await api.request(`/courses/${courseId}/available-students`);
            const allStudents = res.items ?? [];  // теперь это точно массив
            setUsersForEnrollment(allStudents);
        } catch (err) {
            console.error(err);
        }
    };




    useEffect(() => { loadCourses(); }, []);

    const handleCreate = async e => {
        e.preventDefault(); 
        setLoading(true);
        try {
            await api.createCourse({ title: newTitle });
            setNewTitle(''); 
            setShowCreate(false); 
            loadCourses();
        } catch (err) { 
            alert(err.message); 
        } finally { 
            setLoading(false); 
        }
    };

    const handleEnrollClick = (course) => {
        setCurrentCourseId(course.id);
        setShowEnroll(course);
        loadStudentsAvailableForEnrollment(course.id); // load students only when opening modal
    };

    const handleMarkClick = (course) => {
        setCurrentCourseId(course.id);
        setShowMark(course);
        loadStudents(course.id); // load students only when opening modal
    };

    return (
        <div className="dashboard-content">
            <div className="page-header">
                <h1 className="page-title">My Courses</h1>
                <button className="btn-primary" onClick={() => setShowCreate(true)}>
                    <Plus size={16} /> Create Course
                </button>
            </div>

            {courses.length === 0
                ? <EmptyState icon={<BookOpen size={48} />} text="No courses yet" />
                : <div className="courses-grid">
                    {courses.map(c => (
                        <div key={c.id} className="course-card">
                            <div className="course-header"><h3>{c.title}</h3></div>
                            <div className="course-meta"><span>ID: {c.id}</span></div>
                            <div className="course-footer">
                                <button className="btn-secondary" onClick={() => handleEnrollClick(c)}>
                                    <UserPlus size={14} /> Enroll Student
                                </button>
                                <button className="btn-secondary" onClick={() => handleMarkClick(c)}>
                                    <CheckCircle size={14} /> Mark Attendance
                                </button>
                            </div>
                        </div>
                    ))}
                </div>}

            {showCreate && (
                <Modal title="Create Course" onClose={() => setShowCreate(false)}>
                    <form onSubmit={handleCreate}>
                        <div className="form-group">
                            <label>Course Title</label>
                            <input 
                                type="text" 
                                value={newTitle} 
                                onChange={e => setNewTitle(e.target.value)} 
                                required 
                                placeholder="e.g. Introduction to Programming" 
                            />
                        </div>
                        <div className="modal-actions">
                            <button type="button" className="btn-secondary" onClick={() => setShowCreate(false)}>Cancel</button>
                            <button type="submit" className="btn-primary" disabled={loading}>{loading ? 'Creating…' : 'Create'}</button>
                        </div>
                    </form>
                </Modal>
            )}

            {showEnroll && (
                <EnrollModal 
                    course={showEnroll} 
                    students={usersForEnrollment} 
                    onClose={() => { setShowEnroll(null); loadStudents(currentCourseId); }} 
                />
            )}

            {showMark && (
                <MarkAttendanceModal 
                    course={showMark} 
                    students={usersInCourse} 
                    onClose={() => setShowMark(null)} 
                />
            )}
        </div>
    );
}


function TeacherAttendance() {
    const [courses, setCourses] = useState([]);
    const [selected, setSelected] = useState(null);
    const [records, setRecords] = useState([]);

    useEffect(() => { api.getMyCourses().then(setCourses).catch(console.error); }, []);

    const selectCourse = async c => {
        setSelected(c);
        console.log('Selected course:', c.id);
        const data = await api.getCourseAttendance(c.id).catch(() => []);
        setRecords(data);
    };

    return (
        <div className="dashboard-content">
            <h1 className="page-title">Attendance Records</h1>
            <div className="attendance-container">
                <CourseSidebar courses={courses} selected={selected} onSelect={selectCourse} />
                <div className="attendance-list">
                    {!selected
                        ? <EmptyState icon={<Calendar size={48} />} text="Select a course" />
                        : records.length === 0
                            ? <EmptyState icon={<Calendar size={48} />} text="No records yet" />
                            : <AttendanceTable records={records} showStudent={true} />}
                </div>
            </div>
        </div>
    );
}

// ─── ADMIN views ────────────────────────────────────────────────────────────

function AdminHome({ user }) {
    const [stats, setStats] = useState({ users: 0, courses: 0, admins: 0, teachers: 0, students: 0 });

    useEffect(() => {
        Promise.all([api.getUsers(), api.getCourses()]).then(([users, courses]) => {
            setStats({
                users:    users.length,
                courses:  courses.length,
                admins:   users.filter(u => u.role_id === 1).length,
                teachers: users.filter(u => u.role_id === 2).length,
                students: users.filter(u => u.role_id === 3).length,
            });
        }).catch(console.error);
    }, []);

    return (
        <div className="dashboard-content">
            <div className="welcome-section">
                <h1>Welcome back, {user.full_name}!</h1>
                <p className="subtitle">System administration overview</p>
            </div>
            <div className="stats-grid">
                <StatCard icon={<Users />}    value={stats.users}    label="Total Users"    color="#3b6ea8" />
                <StatCard icon={<BookOpen />} value={stats.courses}  label="Total Courses"  color="#5fa8a0" />
                <StatCard icon={<Shield />}   value={stats.admins}   label="Admins"         color="#e74c3c" />
                <StatCard icon={<User />}     value={stats.teachers} label="Teachers"       color="#3b6ea8" />
                <StatCard icon={<GraduationCap />} value={stats.students} label="Students"  color="#5fa8a0" />
            </div>
        </div>
    );
}

function AdminUsers() {
    const [users, setUsers] = useState([]);
    const [showCreate, setShowCreate] = useState(false);
    const [loading, setLoading] = useState(false);
    const [form, setForm] = useState({ email: '', password: '', full_name: '', role_id: 3 });
    const [changingRole, setChangingRole] = useState(null); // { userId, current }

    const load = () => api.getUsers().then(setUsers).catch(console.error);
    useEffect(() => { load(); }, []);

    const set = k => e => setForm(f => ({ ...f, [k]: e.target.value }));

    const handleCreate = async e => {
        e.preventDefault(); setLoading(true);
        try {
            await api.createUser({ ...form, role_id: Number(form.role_id) });
            setShowCreate(false);
            setForm({ email: '', password: '', full_name: '', role_id: 3 });
            load();
        } catch (err) { alert(err.message); }
        finally { setLoading(false); }
    };

    const handleRoleChange = async (userId, newRole) => {
        try {
            await api.updateUserRole(userId, newRole);
            setChangingRole(null);
            load();
        } catch (err) { alert(err.message); }
    };

    const roleMap = { 1: 'admin', 2: 'teacher', 3: 'student' };

    return (
        <div className="dashboard-content">
            <div className="page-header">
                <h1 className="page-title">User Management</h1>
                <button className="btn-primary" onClick={() => setShowCreate(true)}><Plus size={16} /> Create User</button>
            </div>

            <div className="table-card">
                <table className="data-table">
                    <thead>
                        <tr>
                            <th>ID</th><th>Name</th><th>Email</th><th>Role</th><th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {users.map(u => (
                            <tr key={u.id}>
                                <td className="text-muted">#{u.id}</td>
                                <td><strong>{u.full_name}</strong></td>
                                <td className="text-muted">{u.email}</td>
                                <td><RolePill role={roleMap[u.role_id] || 'unknown'} /></td>
                                <td>
                                    {changingRole?.userId === u.id ? (
                                        <div className="role-change-row">
                                            <Select
                                                value={changingRole.newRole}
                                                onChange={v => setChangingRole(r => ({ ...r, newRole: v }))}
                                                options={[
                                                    { value: 'admin', label: 'Admin' },
                                                    { value: 'teacher', label: 'Teacher' },
                                                    { value: 'student', label: 'Student' },
                                                ]}
                                            />
                                            <button className="btn-primary btn-sm" onClick={() => handleRoleChange(u.id, changingRole.newRole)}>Save</button>
                                            <button className="btn-secondary btn-sm" onClick={() => setChangingRole(null)}>Cancel</button>
                                        </div>
                                    ) : (
                                        <button className="btn-secondary btn-sm" onClick={() => setChangingRole({ userId: u.id, newRole: roleMap[u.role_id] || 'student' })}>
                                            <Settings size={13} /> Change Role
                                        </button>
                                    )}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {showCreate && (
                <Modal title="Create User" onClose={() => setShowCreate(false)}>
                    <form onSubmit={handleCreate}>
                        <div className="form-group">
                            <label>Full Name</label>
                            <input type="text" value={form.full_name} onChange={set('full_name')} required placeholder="John Doe" />
                        </div>
                        <div className="form-group">
                            <label>Email</label>
                            <input type="email" value={form.email} onChange={set('email')} required placeholder="john@example.com" />
                        </div>
                        <div className="form-group">
                            <label>Password</label>
                            <input type="password" value={form.password} onChange={set('password')} required placeholder="••••••••" />
                        </div>
                        <div className="form-group">
                            <label>Role</label>
                            <Select value={form.role_id} onChange={v => setForm(f => ({ ...f, role_id: Number(v) }))}
                                options={[
                                    { value: 1, label: 'Admin' },
                                    { value: 2, label: 'Teacher' },
                                    { value: 3, label: 'Student' },
                                ]}
                            />
                        </div>
                        <div className="modal-actions">
                            <button type="button" className="btn-secondary" onClick={() => setShowCreate(false)}>Cancel</button>
                            <button type="submit" className="btn-primary" disabled={loading}>{loading ? 'Creating…' : 'Create User'}</button>
                        </div>
                    </form>
                </Modal>
            )}
        </div>
    );
}

function AdminCourses() {
    const [courses, setCourses] = useState([]);
    const [users, setUsers] = useState([]);
    const [showCreate, setShowCreate] = useState(false);
    const [showEnroll, setShowEnroll] = useState(null);
    const [newTitle, setNewTitle] = useState('');
    const [loading, setLoading] = useState(false);

    const load = () => Promise.all([api.getCourses(), api.getUsers()])
        .then(([c, u]) => { setCourses(c); setUsers(u); }).catch(console.error);

    useEffect(() => { load(); }, []);

    const handleCreate = async e => {
        e.preventDefault(); setLoading(true);
        try {
            await api.createCourse({ title: newTitle });
            setNewTitle(''); setShowCreate(false); load();
        } catch (err) { alert(err.message); }
        finally { setLoading(false); }
    };

    const teacherName = id => {
        const t = users.find(u => u.id === id);
        return t ? t.full_name : `#${id}`;
    };

    return (
        <div className="dashboard-content">
            <div className="page-header">
                <h1 className="page-title">Course Management</h1>
                <button className="btn-primary" onClick={() => setShowCreate(true)}><Plus size={16} /> Create Course</button>
            </div>

            <div className="table-card">
                <table className="data-table">
                    <thead><tr><th>ID</th><th>Title</th><th>Teacher</th><th>Created</th><th>Actions</th></tr></thead>
                    <tbody>
                        {courses.map(c => (
                            <tr key={c.id}>
                                <td className="text-muted">#{c.id}</td>
                                <td><strong>{c.title}</strong></td>
                                <td>{teacherName(c.teacher_id)}</td>
                                <td className="text-muted">{c.created_at ? new Date(c.created_at).toLocaleDateString() : '—'}</td>
                                <td>
                                    <button className="btn-secondary btn-sm" onClick={() => setShowEnroll(c)}>
                                        <UserPlus size={13} /> Enroll Student
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {showCreate && (
                <Modal title="Create Course" onClose={() => setShowCreate(false)}>
                    <form onSubmit={handleCreate}>
                        <div className="form-group">
                            <label>Course Title</label>
                            <input type="text" value={newTitle} onChange={e => setNewTitle(e.target.value)} required placeholder="e.g. Advanced Mathematics" />
                        </div>
                        <div className="modal-actions">
                            <button type="button" className="btn-secondary" onClick={() => setShowCreate(false)}>Cancel</button>
                            <button type="submit" className="btn-primary" disabled={loading}>{loading ? 'Creating…' : 'Create'}</button>
                        </div>
                    </form>
                </Modal>
            )}

            {showEnroll && (
                <EnrollModal course={showEnroll} students={users.filter(u => u.role_id === 3)} onClose={() => { setShowEnroll(null); load(); }} />
            )}
        </div>
    );
}

function AdminAttendance() {
    const [courses, setCourses] = useState([]);
    const [selected, setSelected] = useState(null);
    const [records, setRecords] = useState([]);

    useEffect(() => { api.getCourses().then(setCourses).catch(console.error); }, []);

    const selectCourse = async c => {
        setSelected(c);
        const data = await api.getCourseAttendance(c.id).catch(() => []);
        setRecords(data);
    };

    return (
        <div className="dashboard-content">
            <h1 className="page-title">Attendance Overview</h1>
            <div className="attendance-container">
                <CourseSidebar courses={courses} selected={selected} onSelect={selectCourse} />
                <div className="attendance-list">
                    {!selected
                        ? <EmptyState icon={<Calendar size={48} />} text="Select a course" />
                        : records.length === 0
                            ? <EmptyState icon={<Calendar size={48} />} text="No records yet" />
                            : <AttendanceTable records={records} showStudent={true} />}
                </div>
            </div>
        </div>
    );
}

// ─── Shared modals ──────────────────────────────────────────────────────────

function EnrollModal({ course, students, onClose }) {
    const [studentId, setStudentId] = useState('');
    const [loading, setLoading] = useState(false);
    const [msg, setMsg] = useState('');

    const handleEnroll = async e => {
        e.preventDefault();
        if (!studentId) return;
        setLoading(true); setMsg('');
        try {
            await api.enrollStudent(course.id, Number(studentId));
            setMsg('Student enrolled successfully!');
            setStudentId('');
        } catch (err) { setMsg('Error: ' + err.message); }
        finally { setLoading(false); }
    };

    return (
        <Modal title={`Enroll Student — ${course.title}`} onClose={onClose}>
            <form onSubmit={handleEnroll}>
                <div className="form-group">
                    <label>Student</label>
                    <Select
                        value={studentId}
                        onChange={setStudentId}
                        placeholder="Select a student…"
                        options={students.map(s => ({ value: s.id, label: `${s.full_name} (${s.email})` }))}
                    />
                </div>
                {msg && <div className={`alert ${msg.startsWith('Error') ? 'error' : 'success'}`}>{msg}</div>}
                <div className="modal-actions">
                    <button type="button" className="btn-secondary" onClick={onClose}>Close</button>
                    <button type="submit" className="btn-primary" disabled={loading || !studentId}>{loading ? 'Enrolling…' : 'Enroll'}</button>
                </div>
            </form>
        </Modal>
    );
}

function MarkAttendanceModal({ course, students, onClose }) {
    const today = new Date().toISOString().split('T')[0]; // YYYY-MM-DD
  console.log('MarkAttendanceModal students:', students);

  // store student_id as string to match Select
  const [form, setForm] = useState({
    studentId: '',
    lesson_date: today,
    status: 'present',
    note: '',
  });
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState('');

  const set = k => v => setForm(f => ({ ...f, [k]: v }));

  const handleMark = async e => {
    e.preventDefault();
    if (!form.studentId) return; // nothing selected
    setLoading(true);
    setMsg('');
    // const isoDate = new Date(form.lesson_date + 'T00:00:00Z').toISOString();
    try {
      await api.markAttendance(course.id, {
        student_id: Number(form.studentId), // convert to number for API
        lesson_date: new Date(form.lesson_date + 'T00:00:00Z').toISOString(),
        status: form.status,
        note: form.note,
      });
      setMsg('Attendance marked!');
      setForm(f => ({ ...f, studentId: '' })); // reset selection
    } catch (err) {
      setMsg('Error: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal title={`Mark Attendance — ${course.title}`} onClose={onClose}>
      <form onSubmit={handleMark}>
        <div className="form-group">
            <label>Student</label>
            <Select
                value={form.studentId}
                onChange={set('studentId')}
                placeholder="Select a student…"
                options={students.map(s => ({ value: s.id, label: `${s.full_name} (${s.email})` }))}
            />
        </div>

        <div className="form-group">
          <label>Date</label>
          <input
            type="date"
            value={form.lesson_date}
            onChange={e => set('lesson_date')(e.target.value)}
            required
          />
        </div>

        <div className="form-group">
          <label>Status</label>
          <Select
            value={form.status}
            onChange={set('status')}
            options={[
              { value: 'present', label: '✅ Present' },
              { value: 'absent', label: '❌ Absent' },
              { value: 'late', label: '⏰ Late' },
            ]}
          />
        </div>

        <div className="form-group">
          <label>Note (optional)</label>
          <input
            type="text"
            value={form.note}
            onChange={e => set('note')(e.target.value)}
            placeholder="Any note…"
          />
        </div>

        {msg && (
          <div className={`alert ${msg.startsWith('Error') ? 'error' : 'success'}`}>
            {msg}
          </div>
        )}

        <div className="modal-actions">
          <button type="button" className="btn-secondary" onClick={onClose}>
            Close
          </button>
          <button
            type="submit"
            className="btn-primary"
            disabled={loading || !form.studentId} // use studentId, not student_id
            >
            {loading ? 'Saving…' : 'Mark'}
        </button>

        </div>
      </form>
    </Modal>
  );
}



// ─── Reusable components ────────────────────────────────────────────────────

function StatCard({ icon, value, label, color }) {
    return (
        <div className="stat-card">
            <div className="stat-icon" style={color ? { background: `${color}1a`, color } : {}}>{icon}</div>
            <div className="stat-info"><h3>{value}</h3><p>{label}</p></div>
        </div>
    );
}

function CourseCard({ course }) {
    return (
        <div className="course-card">
            <div className="course-header"><h3>{course.title}</h3></div>
            <div className="course-meta">
                <span>Course ID: {course.id}</span>
                <span>Teacher ID: {course.teacher_id}</span>
            </div>
        </div>
    );
}

function EmptyState({ icon, text }) {
    return (
        <div className="empty-state">
            {icon}
            <p>{text}</p>
        </div>
    );
}

function CourseSidebar({ courses, selected, onSelect }) {
    return (
        <div className="course-selector">
            <h3>Courses</h3>
            {courses.length === 0
                ? <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>No courses</p>
                : courses.map(c => (
                    <div key={c.id}
                        className={`course-item ${selected?.id === c.id ? 'active' : ''}`}
                        onClick={() => onSelect(c)}>
                        <BookOpen size={16} /><span>{c.title}</span>
                    </div>
                ))}
        </div>
    );
}

function AttendanceTable({ records, showStudent }) {
    return (
        <div className="attendance-table">
            <table>
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Status</th>
                        {showStudent && <th>Student ID</th>}
                        <th>Note</th>
                    </tr>
                </thead>
                <tbody>
                    {records.map(r => (
                        <tr key={r.id}>
                            <td>{r.lesson_date}</td>
                            <td><span className={`status-badge ${r.status}`}>{r.status}</span></td>
                            {showStudent && <td>#{r.student_id}</td>}
                            <td>{r.note || '—'}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

// ─── Navigation config per role ─────────────────────────────────────────────

const NAV = {
    student: [
        { id: 'home',       label: 'Dashboard', icon: Home },
        { id: 'courses',    label: 'Courses',   icon: BookOpen },
        { id: 'attendance', label: 'Attendance', icon: Calendar },
        { id: 'profile',    label: 'Profile',   icon: User },
    ],
    teacher: [
        { id: 'home',       label: 'Dashboard', icon: Home },
        { id: 'courses',    label: 'My Courses', icon: BookOpen },
        { id: 'attendance', label: 'Attendance', icon: Calendar },
        { id: 'profile',    label: 'Profile',   icon: User },
    ],
    admin: [
        { id: 'home',       label: 'Dashboard', icon: Home },
        { id: 'users',      label: 'Users',     icon: Users },
        { id: 'courses',    label: 'Courses',   icon: BookOpen },
        { id: 'attendance', label: 'Attendance', icon: Calendar },
        { id: 'profile',    label: 'Profile',   icon: User },
    ],
};

const PAGES = {
    student: {
        home:       u => <StudentHome user={u} />,
        courses:    u => <StudentCourses user={u} />,
        attendance: u => <StudentAttendance user={u} />,
        profile:    u => <ProfilePage user={u} />,
    },
    teacher: {
        home:       u => <TeacherHome user={u} />,
        courses:    u => <TeacherCourses user={u} />,
        attendance: u => <TeacherAttendance user={u} />,
        profile:    u => <ProfilePage user={u} />,
    },
    admin: {
        home:       u => <AdminHome user={u} />,
        users:      u => <AdminUsers user={u} />,
        courses:    u => <AdminCourses user={u} />,
        attendance: u => <AdminAttendance user={u} />,
        profile:    u => <ProfilePage user={u} />,
    },
};

// ─── Dashboard shell ────────────────────────────────────────────────────────

function Dashboard({ user, onLogout }) {
    const role = user.role || 'student';
    const nav   = NAV[role]   || NAV.student;
    const pages = PAGES[role] || PAGES.student;
    const [page, setPage] = useState('home');
    const [sidebarOpen, setSidebarOpen] = useState(true);

    return (
        <div className="dashboard">
            <aside className={`sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
                <div className="sidebar-header">
                    <GraduationCap size={32} />
                    {sidebarOpen && <h2>AITU LMS</h2>}
                </div>
                {sidebarOpen && (
                    <div className="sidebar-role-badge">
                        <RolePill role={role} />
                    </div>
                )}
                <nav className="sidebar-nav">
                    {nav.map(item => {
                        const Icon = item.icon;
                        return (
                            <button key={item.id}
                                className={`nav-item ${page === item.id ? 'active' : ''}`}
                                onClick={() => setPage(item.id)}>
                                <Icon size={20} />
                                {sidebarOpen && <span>{item.label}</span>}
                            </button>
                        );
                    })}
                </nav>
                <div className="sidebar-footer">
                    <button className="nav-item" onClick={onLogout}>
                        <LogOut size={20} />
                        {sidebarOpen && <span>Logout</span>}
                    </button>
                </div>
            </aside>

            <main className="main-content">
                <header className="top-bar">
                    <button className="menu-toggle" onClick={() => setSidebarOpen(o => !o)}>
                        {sidebarOpen ? <X /> : <Menu />}
                    </button>
                    <div className="user-info">
                        <span>{user.full_name}</span>
                        <UserCircle size={32} />
                    </div>
                </header>
                <div className="content-area">
                    {(pages[page] || pages.home)(user)}
                </div>
            </main>
        </div>
    );
}

// ─── App root ───────────────────────────────────────────────────────────────

export default function App() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    const checkAuth = async () => {
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const userData = await api.getProfile();
                setUser(userData);
                setIsAuthenticated(true);
            } catch {
                localStorage.removeItem('token');
            }
        }
        setLoading(false);
    };

    useEffect(() => { checkAuth(); }, []);

    const handleLogout = () => {
        localStorage.removeItem('token');
        setIsAuthenticated(false);
        setUser(null);
    };

    if (loading) return (
        <div className="loading-screen">
            <GraduationCap size={64} className="pulse" />
            <p>Loading…</p>
        </div>
    );

    return (
        <AuthContext.Provider value={{ user, isAuthenticated, logout: handleLogout }}>
            {isAuthenticated && user
                ? <Dashboard user={user} onLogout={handleLogout} />
                : <LoginPage onLogin={checkAuth} />}
        </AuthContext.Provider>
    );
}